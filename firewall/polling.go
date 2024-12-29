package firewall

import (
	"log"
	"sync"
	"time"
)

type PollingFirewallBackend interface {
	AddRule(rule FirewallRule) error
	RemoveRule(rule FirewallRule) error
}

const DefaultPollingInterval = 1 * time.Second

// PollingFirewallAdapter is needed for firewalls that do not support expiration times on rules.
type PollingFirewallAdapter struct {
	interval time.Duration
	rules    map[FirewallRule]time.Time
	mutex    *sync.RWMutex
	backend  PollingFirewallBackend
	done     chan struct{}
	ticker   *time.Ticker
}

func NewPollingFirewallAdapter(backend PollingFirewallBackend) FirewallBackend {
	return &PollingFirewallAdapter{
		interval: DefaultPollingInterval,
		backend:  backend,
		rules:    make(map[FirewallRule]time.Time),
		mutex:    &sync.RWMutex{},
		done:     make(chan struct{}),
	}
}

func (p *PollingFirewallAdapter) Start() error {
	if p.ticker != nil {
		log.Printf("PollingFirewallAdapter already started")
		return nil
	}

	p.ticker = time.NewTicker(p.interval)
	go func() {
		for {
			select {
			case <-p.done:
				return
			case <-time.After(p.interval):
				p.poll()
			}
		}
	}()
	return nil
}

func (p *PollingFirewallAdapter) poll() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	numExpired := 0

	now := time.Now()
	for rule, expiration := range p.rules {
		if expiration.Before(now) {
			log.Printf("Removing expired rule (%v)", rule)
			p.backend.RemoveRule(rule)
			delete(p.rules, rule)

			numExpired++
		}
	}

	if numExpired > 0 {
		log.Printf("Removed %v expired rule(s)", numExpired)
	}
}

func (p *PollingFirewallAdapter) Stop() error {
	if p.ticker == nil {
		log.Printf("PollingFirewallAdapter was not started")
		return nil
	}
	p.ticker.Stop()
	close(p.done)

	// Clean up any remaining rules
	p.mutex.Lock()
	defer p.mutex.Unlock()

	numRules := len(p.rules)

	for rule := range p.rules {
		log.Printf("Cleaning up rule (%v)", rule)
		if err := p.backend.RemoveRule(rule); err != nil {
			log.Printf("Failed to remove rule (%v): %v", rule, err)
		}
		delete(p.rules, rule)
	}

	if numRules > 0 {
		log.Printf("Cleaned up %v rule(s)", numRules)
	}

	return nil
}

func (p *PollingFirewallAdapter) AddRule(rule FirewallRuleWithExpiration) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	current_expiration, ok := p.rules[rule.FirewallRule]
	if ok && rule.Expires.Before(current_expiration) {
		log.Printf("New expiration is before current expiration, ignoring")
		return nil
	}

	if !ok {
		log.Printf("Adding new rule (%v) with expiration %v", rule.FirewallRule, rule.Expires)
		if err := p.backend.AddRule(rule.FirewallRule); err != nil {
			return err
		}
	} else {
		log.Printf("Updating existing rule's (%v) expiration to %v", rule.FirewallRule, rule.Expires)
	}

	p.rules[rule.FirewallRule] = rule.Expires
	return nil
}

func (p *PollingFirewallAdapter) ListActiveRules() ([]FirewallRuleWithExpiration, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	active_rules := []FirewallRuleWithExpiration{}
	for rule, expiration := range p.rules {
		active_rules = append(active_rules, FirewallRuleWithExpiration{
			FirewallRule: rule,
			Expires:      expiration,
		})
	}
	return active_rules, nil
}
