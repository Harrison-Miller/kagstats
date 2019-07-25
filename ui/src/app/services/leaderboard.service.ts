import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { LeaderboardResult } from '../models';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class LeaderboardService {

  constructor(private http: HttpClient) { }

  getBaseLeaderboard(): Observable<LeaderboardResult> {
    return this.http.get<LeaderboardResult>('/api/leaderboard')
  }

  getLeaderboard(board: string): Observable<LeaderboardResult> {
    return this.http.get<LeaderboardResult>(`/api/leaderboard/${board.toLowerCase()}`)
  }
}
