import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { PagedResult, Player, BasicStats } from '../models';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})
export class PlayersService {
  searchText: string = "";

  constructor(private http: HttpClient) {}

  getPlayer(playerId): Observable<BasicStats> {
    return this.http.get<BasicStats>(`/api/players/${playerId}/basic`);
  }

  getPlayerName(playerId): Observable<Player> {
    return this.http.get<Player>(`/api/players/${playerId}`);
  }

  getPlayers(): Observable<PagedResult<Player>> {
    return this.http.get<PagedResult<Player>>('/api/players');
  }

  searchPlayers(search: string): Observable<Player[]> {
    search = search.toLowerCase();
    return this.http.get<Player[]>(`/api/players/search/${search}`);
  }
}
