package main

import (
	_ "embed"
	"encoding/json"
	"net/url"
	"strings"
	"sync"

	"github.com/AdguardTeam/urlfilter"
	"github.com/AdguardTeam/urlfilter/filterlist"
	"github.com/AdguardTeam/urlfilter/rules"
)

//go:embed easylist.txt
var easyList string

//go:embed easyprivacy.txt
var easyPrivacy string

var (
	cosmeticOnce sync.Once
	cosmetic     *urlfilter.CosmeticEngine
)

func cosmeticResources(pageURL string) (css, script string) {
	cosmeticOnce.Do(func() {
		storage, err := filterlist.NewRuleStorage([]filterlist.Interface{
			filterlist.NewString(&filterlist.StringConfig{RulesText: easyList, ID: rules.ListID(1)}),
			filterlist.NewString(&filterlist.StringConfig{RulesText: easyPrivacy, ID: rules.ListID(2)}),
		})
		if err != nil {
			panic(err)
		}
		cosmetic = urlfilter.NewCosmeticEngine(storage)
	})
	parsed, err := url.Parse(pageURL)
	if err != nil || parsed.Hostname() == "" {
		return "", ""
	}
	result := cosmetic.Match(parsed.Hostname(), true, true, true)
	selectors := append(result.ElementHiding.Generic, result.ElementHiding.Specific...)
	var styles strings.Builder
	for _, selector := range selectors {
		if strings.TrimSpace(selector) == "" {
			continue
		}
		styles.WriteString(selector)
		styles.WriteString("{display:none!important}\n")
	}
	return styles.String(), strings.Join(append(result.JS.Generic, result.JS.Specific...), "\n")
}

func injectCSS(css string) string {
	encoded, _ := json.Marshal(css)
	return "(function(){var h=document.head||document.documentElement;if(!h)return;var old=document.getElementById('gf-adblock');if(old)old.remove();var s=document.createElement('style');s.id='gf-adblock';s.textContent=" + string(encoded) + ";h.appendChild(s);})();"
}
