import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { PagedResult, Player, BasicStats, APIPlayerStatus } from '../models';
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

  getStatus(playerName: string) : Observable<APIPlayerStatus> {
    return this.http.get<{playerStatus:APIPlayerStatus}>(`https://api.kag2d.com/v1/player/${playerName}/status`).pipe(
    map(status => {
      var b = status.playerStatus.lastUpdate.split(/\D/);
      status.playerStatus.lastUpdateDate = Date.UTC(+b[0], +b[1]-1, +b[2], +b[3], +b[4], +b[5]);
      return status.playerStatus;
    })
  );
  }
}
