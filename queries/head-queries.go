package queries

import (
	"go.mongodb.org/mongo-driver/bson"
)

func GetHeadQuery(batter string, bowler string, match_type string, group_by string) bson.A {
	query := bson.A{
		bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "innings", Value: bson.D{
					{Key: "$elemMatch", Value: bson.D{
						{Key: "overs", Value: bson.D{
							{Key: "$elemMatch", Value: bson.D{
								{Key: "deliveries", Value: bson.D{
									{Key: "$elemMatch", Value: bson.D{
										{Key: "batter", Value: batter},
										{Key: "bowler", Value: bowler},
									}},
								}},
							}},
						}},
					}},
				}},
			}},
		},
		bson.D{
			{Key: "$addFields", Value: bson.D{
				{Key: "innings", Value: bson.D{
					{Key: "$map", Value: bson.D{
						{Key: "input", Value: "$innings"},
						{Key: "as", Value: "i"},
						{Key: "in", Value: bson.D{
							{Key: "overs", Value: bson.D{
								{Key: "$map", Value: bson.D{
									{Key: "input", Value: "$$i.overs"},
									{Key: "as", Value: "o"},
									{Key: "in", Value: bson.D{
										{Key: "deliveries", Value: bson.D{
											{Key: "$filter", Value: bson.D{
												{Key: "input", Value: "$$o.deliveries"},
												{Key: "as", Value: "b"},
												{Key: "cond", Value: bson.D{
													{Key: "$and", Value: bson.A{
														bson.D{
															{Key: "$eq", Value: bson.A{
																"$$b.batter", batter,
															}},
														},
														bson.D{
															{Key: "$eq", Value: bson.A{
																"$$b.bowler", bowler,
															}},
														},
														bson.D{
															{Key: "$lte", Value: bson.A{
																"$$b.extras.wides", false,
															}},
														},
													}},
												}},
											}},
										}},
									}},
								}},
							}},
						}},
					}},
				}},
			}},
		},
		bson.D{
			{Key: "$addFields", Value: bson.D{
				{Key: "innings", Value: bson.D{
					{Key: "$map", Value: bson.D{
						{Key: "input", Value: "$innings"},
						{Key: "as", Value: "i"},
						{Key: "in", Value: bson.D{
							{Key: "overs", Value: bson.D{
								{Key: "$filter", Value: bson.D{
									{Key: "input", Value: "$$i.overs"},
									{Key: "as", Value: "o"},
									{Key: "cond", Value: bson.D{
										{Key: "$and", Value: bson.A{
											bson.D{
												{Key: "$gt", Value: bson.A{
													bson.D{
														{Key: "$size", Value: "$$o.deliveries"},
													},
													0,
												}},
											},
										}},
									}},
								}},
							}},
						}},
					}},
				}},
			}},
		},
		bson.D{
			{Key: "$addFields", Value: bson.D{
				{Key: "innings", Value: bson.D{
					{Key: "$filter", Value: bson.D{
						{Key: "input", Value: "$innings"},
						{Key: "as", Value: "i"},
						{Key: "cond", Value: bson.D{
							{Key: "$and", Value: bson.A{
								bson.D{
									{Key: "$gt", Value: bson.A{
										bson.D{
											{Key: "$size", Value: "$$i.overs"},
										},
										0,
									}},
								},
							}},
						}},
					}},
				}},
			}},
		},
	}

	if match_type == "IT20" {
		query = append(query, bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "info.match_type", Value: "T20"},
				{Key: "info.team_type", Value: "international"},
			}},
		})
	} else if match_type == "T20" {
		query = append(query, bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "info.match_type", Value: "T20"},
				{Key: "info.team_type", Value: "club"},
			}},
		})
	} else if match_type != "" {
		query = append(query, bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "info.match_type", Value: match_type},
			}},
		})
	}

	if group_by != "" {
		query = append(query, bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: bson.D{
					{Key: "match_type", Value: "$info.match_type"},
					{Key: "group_by", Value: "$info." + group_by},
				}},
				{Key: "matches", Value: bson.D{
					{Key: "$addToSet", Value: bson.D{
						{Key: "_id", Value: "$_id"},
						{Key: "meta", Value: "$meta"},
						{Key: "info", Value: "$info"},
						{Key: "innings", Value: "$innings"},
					}},
				}},
				{Key: "total", Value: bson.D{{Key: "$sum", Value: 1}}},
			}},
		})
	} else {
		query = append(query, bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: bson.D{
					{Key: "match_type", Value: "$info.match_type"},
					{Key: "group_by", Value: "$info.season"},
				}},
				{Key: "matches", Value: bson.D{
					{Key: "$addToSet", Value: bson.D{
						{Key: "_id", Value: "$_id"},
						{Key: "meta", Value: "$meta"},
						{Key: "info", Value: "$info"},
						{Key: "innings", Value: "$innings"},
					}},
				}},
				{Key: "total", Value: bson.D{{Key: "$sum", Value: 1}}},
			}},
		})
	}

	query = append(query, bson.D{
		{Key: "$sort", Value: bson.D{
			{Key: "_id.group_by", Value: 1},
		}},
	})

	return query
}
