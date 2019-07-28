import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { LeaderboardResult } from '../models';
import { BehaviorSubject, Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class LeaderboardService {

  leaderboard = new BehaviorSubject<LeaderboardResult>({
    size: 0,
    leaderboard: []
  });

  constructor(private http: HttpClient) { }

  getBaseLeaderboard(): void {
    this.http.get<LeaderboardResult>('/api/leaderboard').subscribe(board => {
      this.leaderboard.next(board);
    });
  }

  getLeaderboard(board: string): void {
    this.http.get<LeaderboardResult>(`/api/leaderboard/${board.toLowerCase()}`).subscribe(board => {
      this.leaderboard.next(board);
    });
  }
}
