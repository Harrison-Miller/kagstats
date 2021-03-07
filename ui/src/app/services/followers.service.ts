import { Injectable } from '@angular/core';
import {Observable} from 'rxjs';
import {environment} from '../../environments/environment';
import {HttpClient} from '@angular/common/http';
import {BasicStats, FollowingCount} from '../models';

@Injectable({
  providedIn: 'root'
})
export class FollowersService {

  constructor(private http: HttpClient) { }

  getFollowersCount(playerID: number): Observable<FollowingCount> {
    const path = `${environment.apiUrl}/players/${playerID}/followers`;
    return this.http.get<FollowingCount>(path);
  }

  followPlayer(playerID: number): Observable<any> {
    const path = `${environment.apiUrl}/players/${playerID}/follow`;
    return this.http.post<any>(path, {});
  }

  isFollowingPlayer(playerID: number): Observable<FollowingCount> {
    const path = `${environment.apiUrl}/players/${playerID}/follow`;
    return this.http.get<FollowingCount>(path);
  }

  unfollowPlayer(playerID: number): Observable<any> {
    const path = `${environment.apiUrl}/players/${playerID}/unfollow`;
    return this.http.post<any>(path, {});
  }

  getFollowedStats(): Observable<BasicStats[]> {
    const path = `${environment.apiUrl}/following/stats`;
    return this.http.get<BasicStats[]>(path);
  }
}
