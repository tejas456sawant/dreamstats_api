package queries

import (
	"go.mongodb.org/mongo-driver/bson"
)

func FormBatterQuery(name string, limit int64, match_type string) bson.A {
	query := bson.A{}

	if name != "" && limit > 0 {
		query = bson.A{
			bson.D{
				{Key: "$sort", Value: bson.D{
					{Key: "info.dates", Value: -1},
				}},
			},
			bson.D{
				{Key: "$addFields", Value: bson.D{
					{Key: "info.players", Value: bson.D{
						{Key: "$objectToArray", Value: "$info.players"},
					}},
				}},
			},
			bson.D{
				{Key: "$match", Value: bson.D{
					{Key: "info.players", Value: bson.D{
						{Key: "$elemMatch", Value: bson.D{
							{Key: "v", Value: name},
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
																{Key: "$eq",
																	Value: bson.A{"$$b.batter", name},
																},
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
			bson.D{
				{Key: "$match", Value: bson.D{
					{Key: "$expr", Value: bson.D{
						{Key: "$gt", Value: bson.A{
							bson.D{
								{Key: "$size", Value: "$innings"},
							},
							0,
						}},
					}},
				}},
			},
			bson.D{
				{Key: "$limit", Value: limit},
			},
		}
	}

	if name != "" && limit > 0 && match_type != "" {
		query = bson.A{
			bson.D{
				{Key: "$sort", Value: bson.D{
					{Key: "info.dates", Value: -1},
				}},
			},
			bson.D{
				{Key: "$addFields", Value: bson.D{
					{Key: "info.players", Value: bson.D{
						{Key: "$objectToArray", Value: "$info.players"},
					}},
				}},
			},
			bson.D{
				{Key: "$match", Value: bson.D{
					{Key: "info.players", Value: bson.D{
						{Key: "$elemMatch", Value: bson.D{
							{Key: "v", Value: name},
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
																{Key: "$eq",
																	Value: bson.A{"$$b.batter", name},
																},
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
			bson.D{
				{Key: "$match", Value: bson.D{
					{Key: "$expr", Value: bson.D{
						{Key: "$gt", Value: bson.A{
							bson.D{
								{Key: "$size", Value: "$innings"},
							},
							0,
						}},
					}},
					{Key: "info.match_type", Value: match_type},
				}},
			},
			bson.D{
				{Key: "$limit", Value: limit},
			},
		}
	}

	return query
}

func FormBowlerQuery(name string, limit int64, match_type string) bson.A {
	query := bson.A{}

	if name != "" && limit > 0 {
		query = bson.A{
			bson.D{
				{Key: "$sort", Value: bson.D{
					{Key: "info.dates", Value: -1},
				}},
			},
			bson.D{
				{Key: "$addFields", Value: bson.D{
					{Key: "info.players", Value: bson.D{
						{Key: "$objectToArray", Value: "$info.players"},
					}},
				}},
			},
			bson.D{
				{Key: "$match", Value: bson.D{
					{Key: "info.players", Value: bson.D{
						{Key: "$elemMatch", Value: bson.D{
							{Key: "v", Value: name},
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
																{Key: "$eq",
																	Value: bson.A{"$$b.bowler", name},
																},
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
			bson.D{
				{Key: "$match", Value: bson.D{
					{Key: "$expr", Value: bson.D{
						{Key: "$gt", Value: bson.A{
							bson.D{
								{Key: "$size", Value: "$innings"},
							},
							0,
						}},
					}},
				}},
			},
			bson.D{
				{Key: "$limit", Value: limit},
			},
		}
	}

	if name != "" && limit > 0 && match_type != "" {
		query = bson.A{
			bson.D{
				{Key: "$sort", Value: bson.D{
					{Key: "info.dates", Value: -1},
				}},
			},
			bson.D{
				{Key: "$addFields", Value: bson.D{
					{Key: "info.players", Value: bson.D{
						{Key: "$objectToArray", Value: "$info.players"},
					}},
				}},
			},
			bson.D{
				{Key: "$match", Value: bson.D{
					{Key: "info.players", Value: bson.D{
						{Key: "$elemMatch", Value: bson.D{
							{Key: "v", Value: name},
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
																{Key: "$eq",
																	Value: bson.A{"$$b.bowler", name},
																},
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
			bson.D{
				{Key: "$match", Value: bson.D{
					{Key: "$expr", Value: bson.D{
						{Key: "$gt", Value: bson.A{
							bson.D{
								{Key: "$size", Value: "$innings"},
							},
							0,
						}},
					}},
					{Key: "info.match_type", Value: match_type},
				}},
			},
			bson.D{
				{Key: "$limit", Value: limit},
			},
		}
	}

	return query
}
