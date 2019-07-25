export interface Event {
  id: Number;
  type: String;
  time: Number;
  serverId: Number;
  player: Player;
}

export interface Kill {
  id: Number;
  killerClass: String;
  victimClass: String;
  hitter: Number;
  time: Number;
  serverId: Number;
  teamKill: Number;
  player: Player;
  killer: Player;
}

export interface Player {
  id: Number;
  username: String;
  characterName: String;
  clanTag: String;
}

export interface Server {
  id: Number;
  name: String;
  description: String;
  gameMode: String;
  tags: String;
}

export interface PagedResult<T> {
  limit: Number;
  start: Number;
  size: Number;
  next: Number;
  values: T[];
}
