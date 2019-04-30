#define SERVER_ONLY

void onPlayerDie(CRules@ this, CPlayer@ victim, CPlayer@ killer, u8 customData)
{
	if(sv_tcpr && victim !is null)
	{
		array<string> jsonProps = {
			JSONString("VictimUsername", victim.getUsername()),
			JSONString("VictimCharacterName", victim.getCharacterName()),
			JSONString("VictimClantag", victim.getClantag()),
			JSONString("VictimClass", victim.lastBlobName),
			JSONInt("Hitter", customData)
		};

		if(killer !is null)
		{
			jsonProps.push_back(JSONString("KillerUsername", killer.getUsername()));
			jsonProps.push_back(JSONString("KillerCharacterName", killer.getCharacterName()));
			jsonProps.push_back(JSONString("KillerClantag", killer.getClantag()));
			jsonProps.push_back(JSONString("KillerClass", killer.lastBlobName));
		}

		string object = "{" + join(jsonProps, ",") + "}";
		tcpr("*STATS " + object);
	}
}

string JSONProp(string name)
{
	return "\"" + name + "\":";
}

string JSONString(string name, string value)
{
	return JSONProp(name) + "\"" + value + "\"";

}

string JSONNull(string name)
{
	return JSONProp(name) + "null";

}

string JSONInt(string name, int value)
{
	return JSONProp(name) + value;

}
