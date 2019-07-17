array<string> players;
for(int i = 0; i < getPlayersCount(); i++) {
	CPlayer@ p = getPlayer(i);
	string playerString = "{\"username\":\"" + p.getUsername() + "\",";
	playerString += "\"charactername\":\"" + p.getCharacterName() + "\",";
	if (p.exists("stats_id")) {
		playerString += "\"ID\":"+ p.get_u32("stats_id") + ",";
	}
	playerString += "\"clantag\":\"" + p.getClantag() + "\"}";
	players.push_back(playerString);
}

if(players.length() != 0) {
	string object="PlayerList " + "[" + join(players, ",") + "]";
	tcpr(object);
}
