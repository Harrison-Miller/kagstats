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
  victim: Player;
  killer: Player;
}

export interface Player {
  id: number;
  username: string;
  characterName: string;
  clanTag: string;
  avatar: string;
  oldGold: boolean;
  registered: string;
  role: number;
  tier: number;
  lastEvent: string;
  lastEventTime: number;
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
