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

export interface MonthlyStats extends BasicStats {
  year: number;
  month: number;
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

export interface HittersStats {
  player: Player;

  died: number;
  crushing: number;
  fall: number;
  water: number;
  water_stun: number;
  water_stun_force: number;
  drowning: number;
  fire: number;
  burn: number;
  flying: number;
  stomp: number;
  suicide: number;
  bite: number;
  pickaxe: number;
  sword: number;
  shield: number;
  bomb: number;
  stab: number;
  arrow: number;
  bomb_arrow: number;
  ballista_bolt: number;
  catapult_stones: number;
  catapult_boulder: number;
  boulder: number;
  ram: number;
  explosion: number;
  keg: number;
  mine: number;
  mine_special: number;
  spikes: number;
  saw: number;
  drill: number;
  muscles: number;
  sudden_gib: number;
}

export interface MonthlyHittersStats extends HittersStats {
  year: number;
  month: number;
}

export interface Captures {
  playerID: number;
  captures: number;
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
  charactername: string;
  clantag: string;
  clanID: number;
  clanInfo: ClanInfo;
  joinedClan: number;

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

}

export interface APIServer {
  currentPlayers: number;
  maxPlayers: number;
  firstSeen: string;
  lastUpdate: string;
  description: string;
  serverName: string;
}

export interface APIPlayerStatus {
  action: number;
  lastUpdate: string;
  lastUpdateDate: number;
  server: APIPlayerServer;
  apiServer: APIServer;
}

export interface APIPlayerServer {
  serverIPv4Address: string;
  serverPort: string;
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
  gamemode: string;
}

export interface GithubTreeEntry {
  path: string;
}

export interface GithubTree {
  tree: GithubTreeEntry[];
}

export interface PlayerClaims {
  playerID: number;
  username: string;
  avatar: string;
  clanID: number;
  bannedFromMakingClans: boolean;
}

export interface LoginResp {
  token: string;
}

export interface ClanInfo {
  id: number;
  name: string;
  createdAt: number;
  leaderID: number;

  membersCount: number;
  leader: Player;
}

export interface ClanInvite {
  clanID: number;
  playerID: string;
  time: number;

  player: Player;
  clan: ClanInfo;
  leader: Player;
}

export interface FollowingCount {
  followingCount: number;
}
