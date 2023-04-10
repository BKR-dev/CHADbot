package responses

// InsultResponse
type InsultingResponse struct {
	topic     string   // name of the InsultResponse
	keywords  []string // slice of keywords that yield a respone
	responses []string // slice of spicy responses
}

// Insults
type Insults struct {
	personalAttack      []string // you sodding tiktak
	personalTitleAttack []string // insults you for/with your title
}

// stores all responses and insults for easy addition by users
// returns allResponses type InsultingResponse, Insult struct
// and slice of strings of all keywords of all responses for matching
func ProvideResponsesAndInsults() ([]InsultingResponse, Insults, []string) {
	var trigger []string

	// insults to add to the response
	insults := Insults{
		personalAttack: []string{ // add your inselt here
			"you sodding tiktak",
			"absolute muppet",
			"complete retard",
			"cum guzzling lunatic",
			"beyond meat enjoyer",
			"incel cuck",
			"vegan soy brainlet",
			"just shut the fuck up fag",
			"bubble tea enjoyer",
			"first grade degenerate"},
		personalTitleAttack: []string{
			"%s is your title!? Fitting...", // add insult for a users title here
			"%s, yeah right, my ass!"},
	}

	// all responses per topic

	// name of the Object
	weebAnime := InsultingResponse{
		topic:    "weebAnime",               // name or topic of the Object
		keywords: []string{"anime", "weeb"}, // keywords that are used as trigger for the response
		responses: []string{
			"weebs are retarded and anime is trash", // actual responses
			"2D women are for brainlets and retards",
			"anime is cringe and fake, go and touch some grass"},
	}

	innawoods := InsultingResponse{
		topic:    "innawoods",
		keywords: []string{"innawoods", "inna woods"},
		responses: []string{
			"forever alone in the woods",
			"pissing in jars to keep some company",
			"getting buttfucked by the local wendingo",
			"starving in the cold is better than buying starbucks soy caramel faggoccino"},
	}

	nineteen11 := InsultingResponse{
		topic:    "1911",
		keywords: []string{"1911"},
		responses: []string{
			"two world wars",
			"chinese made glocks are more versatile"},
	}

	shillCrypto := InsultingResponse{
		topic:    "shill",
		keywords: []string{"shill", "crypto", "bitcoin", "etherium", "NFT"},
		responses: []string{
			"Thanks to Coinbase resiliency and UFC NFTs, crypto is now linked directly to my Wells Fargo account :chris_party:",
			"My NFT ETF i just shorted got me the platinum AMEX so i get paid for spending money i dont have :think_about_it:",
			"bro invest in my hyper value adding NFT web4.2 finTech renewable ecommerce gaming startup, bro"},
	}

	diet := InsultingResponse{
		topic:    "diet",
		keywords: []string{"vegan", "keto"},
		responses: []string{
			"just eat some real food and stop being a cunt",
			"stop pretending you are not being brainwashed by some cucked incels to buy their supplements"},
	}

	opSys := InsultingResponse{
		topic:    "opSys",
		keywords: []string{"linux", "macos", "windows"},
		responses: []string{
			"stop being such a poor and use a real OS",
			"you got your programming socks already?"},
	}

	feds := InsultingResponse{
		topic:    "feds",
		keywords: []string{"fed", "FBI", "CIA", "ATF", "AFT"},
		responses: []string{
			"We've had enough, time to blow this fucker up!",
			"Wow you're so cool! Go, commit a crime :hugs:"},
	}

	mockingParantheses := InsultingResponse{
		topic:     "mockingParantheses",
		keywords:  []string{"(((", ")))"},
		responses: []string{},
	}

	vidja := InsultingResponse{
		topic:    "videogames",
		keywords: []string{"gaming", "gayming", "vidja", "videogames"},
		responses: []string{
			"grow up already, you are not a child anymore, just let it go",
			"you played that!? wow amazing, wayyyyy better than watching netflix to numb your brain right!?"},
	}
	// responses slice of structs
	allResponses := []InsultingResponse{ // if you added a new topic, add it here so its being used
		weebAnime, innawoods, nineteen11, shillCrypto, diet, opSys, feds, mockingParantheses, vidja}

	for _, r := range allResponses {
		trigger = append(trigger, r.keywords...)
	}

	return allResponses, insults, trigger
}
