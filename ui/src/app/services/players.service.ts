import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { PagedResult, Player, APIPlayerStatus, APIPlayerInfo } from '../models';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})
export class PlayersService {
  constructor(private http: HttpClient) {}

  getPlayers(): Observable<PagedResult<Player>> {
    return this.http.get<PagedResult<Player>>('/api/players');
  }

  searchPlayers(search: string): Observable<Player[]> {
    search = search.toLowerCase();
    return this.http.get<Player[]>(`/api/players/search/${search}`);
  }

  getPlayerAvatar(username: string): Observable<string> {
    let path = `https://api.kag2d.com/v1/player/${username}/avatar`
    return this.http.get<{large:string}>(path).pipe(
      map(avatar => avatar.large)
    );
  }

  getAPIPlayer(username: string): Observable<{playerInfo:APIPlayerInfo, playerStatus:APIPlayerStatus}> {
    let path = `https://api.kag2d.com/v1/player/${username}`;
    return this.http.get<{playerInfo:APIPlayerInfo, playerStatus:APIPlayerStatus}>(path);
  }
}
