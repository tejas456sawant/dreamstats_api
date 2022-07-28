package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Meta struct {
	DataVersion string `bson:"data_version" json:"data_version"`
	Created     string `bson:"created" json:"created"`
	Revision    int    `bson:"revision" json:"revision"`
}

type Info struct {
	BallsPerOver int `bson:"balls_per_over" json:"balls_per_over"`
	BowlOut      []struct {
		Bowler  string `bson:"bowler" json:"bowler"`
		Outcome string `bson:"outcome" json:"outcome"`
	} `bson:"bowl_out" json:"bowl_out"`
	City  string   `bson:"city" json:"city"`
	Dates []string `bson:"dates" json:"dates"`
	Event struct {
		Name        string `bson:"name" json:"name"`
		MatchNumber int    `bson:"match_number" json:"match_number"`
		Group       string `bson:"group" json:"group"`
		Stage       string `bson:"stage" json:"stage"`
	} `bson:"event" json:"event"`
	Gender          string   `bson:"gender" json:"gender"`
	MatchType       string   `bson:"match_type" json:"match_type"`
	MatchTypeNumber int      `bson:"match_type_number" json:"match_type_number"`
	Missing         []string `bson:"missing" json:"missing"`
	Officials       struct {
		MatchReferees  []string `bson:"match_referees" json:"match_referees"`
		ReserveUmpires []string `bson:"reserve_umpires" json:"reserve_umpires"`
		TvUmpires      []string `bson:"tv_umpires" json:"tv_umpires"`
		Umpires        []string `bson:"umpires" json:"umpires"`
	} `bson:"officials" json:"officials"`
	Outcome struct {
		By struct {
			Innings int `bson:"innings" json:"innings"`
			Runs    int `bson:"runs" json:"runs"`
			Wickets int `bson:"wickets" json:"wickets"`
		} `bson:"by" json:"by"`
		BowlOut    string `bson:"bowl_out" json:"bowl_out"`
		Eliminator string `bson:"eliminator" json:"eliminator"`
		Method     string `bson:"method" json:"method"`
		Result     string `bson:"result" json:"result"`
		Winner     string `bson:"winner" json:"winner"`
	} `bson:"outcome" json:"outcome"`
	Overs         int      `bson:"overs" json:"overs"`
	PlayerOfMatch []string `bson:"player_of_match" json:"player_of_match"`
	Players       struct{} `bson:"players" json:"players"`
	Registry      struct {
		People struct{} `bson:"people" json:"people"`
	} `bson:"registry" json:"registry"`
	Season    string   `bson:"season" json:"season"`
	Supersubs struct{} `bson:"supersubs" json:"supersubs"`
	TeamType  string   `bson:"team_type" json:"team_type"`
	Teams     []string `bson:"teams" json:"teams"`
	Toss      struct {
		Uncontested bool   `bson:"uncontested" json:"uncontested"`
		Decision    string `bson:"decision" json:"decision"`
		Winner      string `bson:"winner" json:"winner"`
	} `bson:"toss" json:"toss"`
	Venue string `bson:"venue" json:"venue"`
}

type Innings struct {
	Team  string `bson:"team" json:"team"`
	Overs []struct {
		Over       int `bson:"over" json:"over"`
		Deliveries []struct {
			Batter string `bson:"batter" json:"batter"`
			Bowler string `bson:"bowler" json:"bowler"`
			Extras struct {
				Byes    int `bson:"byes" json:"byes"`
				Legbyes int `bson:"legbyes" json:"legbyes"`
				Noballs int `bson:"noballs" json:"noballs"`
				Penalty int `bson:"penalty" json:"penalty"`
				Wides   int `bson:"wides" json:"wides"`
			} `bson:"extras" json:"extras"`
			NonStriker   string `bson:"non_striker" json:"non_striker"`
			Replacements struct {
				Match []struct {
					In     string `bson:"in" json:"in"`
					Out    string `bson:"out" json:"out"`
					Reason string `bson:"reason" json:"reason"`
					Team   string `bson:"team" json:"team"`
				} `bson:"match" json:"match"`
				Role []struct {
					In     string `bson:"in" json:"in"`
					Out    string `bson:"out" json:"out"`
					Reason string `bson:"reason" json:"reason"`
					Role   string `bson:"role" json:"role"`
				} `bson:"role" json:"role"`
			} `bson:"replacements" json:"replacements"`
			Review struct {
				Batter      string `bson:"batter" json:"batter"`
				By          string `bson:"by" json:"by"`
				Decision    string `bson:"decision" json:"decision"`
				Umpire      string `bson:"umpire" json:"umpire"`
				UmpiresCall bool   `bson:"umpires_call" json:"umpires_call"`
			} `bson:"review" json:"review"`
			Runs struct {
				Batsman     int  `bson:"batsman" json:"batsman"`
				Extras      int  `bson:"extras" json:"extras"`
				NonBoundary bool `bson:"non_boundary" json:"non_boundary"`
				Total       int  `bson:"total" json:"total"`
			} `bson:"runs" json:"runs"`
			Wickets []struct {
				Fielders []struct {
					Name string `bson:"name" json:"name"`
				} `bson:"fielders" json:"fielders"`
				Kind       string `bson:"kind" json:"kind"`
				Player_out string `bson:"player_out" json:"player_out"`
			} `bson:"wickets" json:"wickets"`
		} `bson:"deliveries" json:"deliveries"`
	} `bson:"overs" json:"overs"`
	AbsentHurt  []string `bson:"absent_hurt" json:"absent_hurt"`
	PenaltyRuns struct {
		Pre  int `bson:"pre" json:"pre"`
		Post int `bson:"post" json:"post"`
	} `bson:"penalty_runs" json:"penalty_runs"`
	Declared   bool `bson:"declared" json:"declared"`
	Forfeited  bool `bson:"forfeited" json:"forfeited"`
	Powerplays []struct {
		From float64 `bson:"from" json:"from"`
		To   float64 `bson:"to" json:"to"`
		Type string  `bson:"type" json:"type"`
	} `bson:"powerplays" json:"powerplays"`
	// https://cricsheet.org/format/json/#miscounted_overs
	MiscountedOvers struct{} `bson:"miscounted_overs" json:"miscounted_overs"`
	Target          struct {
		Overs int `bson:"overs" json:"overs"`
		Runs  int `bson:"runs" json:"runs"`
	} `bson:"target" json:"target"`
	SuperOver bool `bson:"super_over" json:"super_over"`
}

type Match struct {
	ID      primitive.ObjectID `bson:"_id" json:"_id"`
	Meta    Meta               `bson:"meta" json:"meta"`
	Info    Info               `bson:"info" json:"info"`
	Innings []Innings          `bson:"innings" json:"innings"`
}
