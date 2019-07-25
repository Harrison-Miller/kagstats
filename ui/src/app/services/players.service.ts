import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { PagedResult, Player } from '../models';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class PlayersService {
  constructor(public http: HttpClient) {}

  getPlayers(): Observable<PagedResult<Player>> {
    return this.http.get<PagedResult<Player>>('/api/players');
  }
}
