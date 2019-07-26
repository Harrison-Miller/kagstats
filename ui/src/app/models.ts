export interface BasicStats {
  player: Player;
  suicides: number;
  teamKills: number;
  archerKills: number;
  archerDeaths: number;
  builderKills: number;
  builderDeaths: number;
  knightKills: number;
  knightDeaths: number;
  otherKills: number;
  otherDeaths: number;
}

export interface Event {
  id: number;
  type: string;
  time: number;
  serverId: number;
  player: Player;
}

export interface Kill {
  id: number;
  killerClass: string;
  victimClass: string;
  hitter: number;
  time: number;
  serverId: number;
  teamKill: number;
  player: Player;
  killer: Player;
}

export interface APIPlayerInfo {
  old_gold: boolean;
  registered: string;
  role: number;
}

export interface APIPlayerStatus {
  action: number;
  lastUpdate: string;
}

export interface Player {
  id: number;
  username: string;
  characterName: string;
  clanTag: string;
  avatar: string;
  info: APIPlayerInfo;
  status: APIPlayerStatus;
}

export interface APIServer {
  currentPlayers: number;
  maxPlayers: number;
  firstSeen: string;
  lastUpdate: string;
}

export interface Server {
  id: number;
  name: string;
  description: string;
  gameMode: string;
  address: string;
  port: string;
  tags: string;
  status: APIServer;
}

export interface LeaderboardResult {
  size: number;
  leaderboard: BasicStats[];
}

export interface PagedResult<T> {
  limit: number;
  start: number;
  size: number;
  next: number;
  values: T[];
}
