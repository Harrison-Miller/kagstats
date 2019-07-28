import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { PagedResult, Player } from '../models';
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
}
