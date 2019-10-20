import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { PagedResult, Player, BasicStats } from '../models';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { environment } from '../../environments/environment';

@Injectable({
  providedIn: 'root'
})
export class PlayersService {
  searchText: string = "";

  constructor(private http: HttpClient) {}

  getPlayer(playerId): Observable<BasicStats> {
    return this.http.get<BasicStats>(`${environment.apiUrl}/players/${playerId}/basic`);
  }

  getPlayerName(playerId): Observable<Player> {
    return this.http.get<Player>(`${environment.apiUrl}/players/${playerId}`);
  }

  getPlayers(): Observable<PagedResult<Player>> {
    return this.http.get<PagedResult<Player>>(`${environment.apiUrl}/players`);
  }

  searchPlayers(search: string): Observable<Player[]> {
    search = search.toLowerCase();
    return this.http.get<Player[]>(`${environment.apiUrl}/players/search/${search}`);
  }
}
