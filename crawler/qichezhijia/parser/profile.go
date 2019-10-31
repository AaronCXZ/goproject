package parser

import (
	"crawler/engine"
	"crawler/model"
	"regexp"
	"strconv"
)

var (
	styleRe        = regexp.MustCompile(`<span class="font-arial">([^<]+)</span>`)
	cityRe         = regexp.MustCompile(`<dd class="c333">([^&]+) &nbsp;<span class="js-countryname" data-val="2067510,36242"></span></dd>`)
	dealerRe       = regexp.MustCompile(`<a href="https://dealer.autohome.com.cn/[a-z0-9]+" class="js-dearname" data-val="2067510,36242" data-evalid="2258012" target="_blank">([^<]+)</a>`)
	timerRe        = regexp.MustCompile(`<dd class="font-arial bg-blue">([0-9]+)年([0-9]+)月</dd>`)
	priceRe        = regexp.MustCompile(`<dd class="font-arial bg-blue">([^<]+)<span class="c999">&nbsp;万元</span></dd>`)
	fuelRe         = regexp.MustCompile(`<p>([^<]+)<span class="c999">&nbsp;升/百公里</span></p>`)
	mileageRe      = regexp.MustCompile(`<p>([^<]+)<span class="c999">&nbsp;公里</span></p>`)
	spaceRe        = regexp.MustCompile(`<dl class="choose-dl"><dt>空间.*c333">([^<]+)</span></dd></dl>`)
	powerRe        = regexp.MustCompile(`<dl class="choose-dl"><dt>动力.*c333">([^<]+)</span></dd></dl>`)
	manipulationRe = regexp.MustCompile(`<dl class="choose-dl"><dt>操控.*c333">([^<]+)</span></dd></dl>`)
	consumptionRe  = regexp.MustCompile(`<dl class="choose-dl"><dt>油耗.*c333">([^<]+)</span></dd></dl>`)
	comfortRe      = regexp.MustCompile(`<dl class="choose-dl"><dt>舒适性.*c333">([^<]+)</span></dd></dl>`)
	exteriorRe     = regexp.MustCompile(`<dl class="choose-dl"><dt>外观.*c333">([^<]+)</span></dd></dl>`)
	interiorRe     = regexp.MustCompile(`<dl class="choose-dl"><dt>内饰.*c333">([^<]+)</span></dd></dl>`)
	costRe         = regexp.MustCompile(`<dl class="choose-dl"><dt>性价比.*c333">([^<]+)</span></dd></dl>`)
)

func ParseProfile(contents []byte, name string, url string) engine.ParseResult {
	profile := model.Profile{}

	profile.Name = name
	profile.Style = extractString(contents, styleRe)
	profile.City = extractString(contents, cityRe)
	profile.Dealer = extractString(contents, dealerRe)
	profile.Timer = extractString(contents, timerRe)
	profile.Price = extractString(contents, priceRe)
	profile.Fuel = extractString(contents, fuelRe)
	mileage, err := strconv.Atoi(extractString(contents, mileageRe))
	if err == nil {
		profile.Mileage = mileage
	}
	space, err := strconv.Atoi(extractString(contents, spaceRe))
	if err == nil {
		profile.Space = space
	}
	power, err := strconv.Atoi(extractString(contents, powerRe))
	if err == nil {
		profile.Power = power
	}
	manipulation, err := strconv.Atoi(extractString(contents, manipulationRe))
	if err == nil {
		profile.Manipulation = manipulation
	}
	consumption, err := strconv.Atoi(extractString(contents, consumptionRe))
	if err == nil {
		profile.Consumption = consumption
	}
	comfort, err := strconv.Atoi(extractString(contents, comfortRe))
	if err == nil {
		profile.Comfort = comfort
	}
	exterior, err := strconv.Atoi(extractString(contents, exteriorRe))
	if err == nil {
		profile.Exterior = exterior
	}
	interior, err := strconv.Atoi(extractString(contents, interiorRe))
	if err == nil {
		profile.Interior = interior
	}
	cost, err := strconv.Atoi(extractString(contents, costRe))
	if err == nil {
		profile.Cost = cost
	}
	result := engine.ParseResult{Items: []engine.Item{{
		Url:     url,
		Type:    "",
		Id:      "",
		Payload: profile,
	}}}
	return result
}

func extractString(contents []byte, re *regexp.Regexp) string {
	match := re.FindSubmatch(contents)
	if len(match) >= 2 {
		return string(match[1])
	}
	return ""
}
