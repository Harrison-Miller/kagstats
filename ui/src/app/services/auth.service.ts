import { Injectable } from '@angular/core';
import {BehaviorSubject, Observable} from 'rxjs';
import {LoginResp, PlayerClaims} from '../models';
import {HttpClient} from '@angular/common/http';
import { environment } from '../../environments/environment';
import * as Cookies from 'js-cookie';
import {share} from "rxjs/operators";

const tokenName = 'KAGSTATS_TOKEN';

@Injectable({
  providedIn: 'root'
})
export class AuthService {

  playerClaims = new BehaviorSubject<PlayerClaims>(null);

  constructor(private http: HttpClient) {
    this.getClaims();
  }

  login(username: string, token: string): Observable<LoginResp> {
    const path = `${environment.apiUrl}/login`;
    const obs = this.http.post<LoginResp>(path, {
      username,
      token,
    }, {withCredentials: true}).pipe(share<LoginResp>());
    obs.subscribe( resp => {
      this.playerClaims.next(null);
      Cookies.remove(tokenName);
      Cookies.set(tokenName, resp.token, { expires: 365 });
      this.getClaims();

    });
    return obs;
  }

  validate(): Observable<LoginResp> {
    const path = `${environment.apiUrl}/validate`;
    return this.http.get<any>(path);
  }

  getClaims(): void {
    // skip if we already have playerClaims
    if (this.playerClaims.getValue() != null) {
      return;
    }

    // check if we have a cookie and validate it
    let token = Cookies.get(tokenName);
    if (typeof token !== 'undefined' && token !== '') {
      this.validate().subscribe( resp => {
        Cookies.set(tokenName, resp.token, { expires: 365 });
        token = resp.token;
        const parts = token.split('.');
        if (parts.length === 3) {
          const payload = JSON.parse(atob(parts[1]));
          this.playerClaims.next({
            playerID: payload.playerID,
            username: payload.username,
            avatar: payload.avatar,
            clanID: payload.clanID,
            bannedFromMakingClans: payload.bannedFromMakingClans,
          });
        }
      });
    }

    return;
  }

  revalidate(): void {
    this.playerClaims.next(null);
    this.getClaims();
  }

  logout(): void {
    this.playerClaims.next(null);
    Cookies.remove(tokenName);
  }

}
