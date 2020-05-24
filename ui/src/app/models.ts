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
  totalKills: number;
  totalDeaths: number;
}

export interface Status {
  players: number;
  kills: number;
  servers: number;
  version: string;
}

export interface Nemesis {
  nemesis: Player;
  deaths: number;
}

export interface Hitter {
  hitter: number;
  kills: number;
}

export interface Captures {
  playerID: number;
  captures: number;
}

export interface Event {
  id: number;
  type: string;
  time: number;
  serverId: number;
  playerId: number;
}

export interface Kill {
  id: number;
  killerClass: string;
  victimClass: string;
  hitter: number;
  time: number;
  serverId: number;
  teamKill: boolean;
  victim: Player;
  killer: Player;
  server: Server;
}

export interface Player {
  id: number;
  username: string;
  characterName: string;
  clanTag: string;

  //kag2d.com api
  avatar: string;
  oldGold: boolean;
  registered: string;
  role: number;
  tier: number;

  //accolades
  gold: number;
  silver: number;
  bronze: number;
  participation: number;
  github: boolean;
  community: boolean;
  mapmaker: boolean;
  moderation: boolean;

  lastEvent: Event;


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
  status: boolean;
  APIStatus: APIServer;
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

export interface MapBasics {
  mapName: string;
  average: number;
  stddev: number;
  matches: number;
  ballots: number;
  votes: number;
  wins: number;

  // used for UI only
  image: string;
  percent: number;
}

export interface GithubTreeEntry {
  path: string;
}

export interface GithubTree {
  tree: GithubTreeEntry[];
}