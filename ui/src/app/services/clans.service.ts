import { Injectable } from '@angular/core';
import {Observable} from "rxjs";
import {HttpClient} from "@angular/common/http";
import { environment } from '../../environments/environment';
import {ClanInfo, ClanInvite, Player} from "../models";

@Injectable({
  providedIn: 'root'
})
export class ClansService {

  constructor(private http: HttpClient) { }

  register(clanname: string): Observable<ClanInfo> {
    const path = `${environment.apiUrl}/clans/register`;
    return this.http.post<ClanInfo>(path, {
      name: clanname
    });
  }

  getClan(clanID: number): Observable<ClanInfo> {
    const path = `${environment.apiUrl}/clans/${clanID}`;
    return this.http.get<ClanInfo>(path);
  }

  disband(clanID: number): Observable<any> {
    const path = `${environment.apiUrl}/clans/${clanID}/disband`;
    return this.http.post<any>(path, {});
  }

  invite(clanID: number, username: string): Observable<any> {
    const path = `${environment.apiUrl}/clans/${clanID}/invite`;
    return this.http.post<any>(path, {
      username,
    });
  }

  getClanInvites(clanID: number): Observable<ClanInvite[]> {
    const path = `${environment.apiUrl}/clans/${clanID}/invites`;
    return this.http.get<ClanInvite[]>(path);
  }

  cancelClanInvite(clanID: number, playerID: number): Observable<any> {
    const path = `${environment.apiUrl}/clans/${clanID}/invites/cancel/${playerID}`;
    return this.http.post<any>(path, {});
  }

  getMyInvites(): Observable<ClanInvite[]> {
    const path = `${environment.apiUrl}/clans/invites`;
    return this.http.get<ClanInvite[]>(path);
  }

  declineClanInvite(clanID: number): Observable<any> {
    const path = `${environment.apiUrl}/clans/${clanID}/decline`;
    return this.http.post<any>(path, {});
  }

  acceptClanInvite(clanID: number): Observable<any> {
    const path = `${environment.apiUrl}/clans/${clanID}/accept`;
    return this.http.post<any>(path, {});
  }

  getMembers(clanID: number): Observable<Player[]> {
    const path = `${environment.apiUrl}/clans/${clanID}/members`;
    return this.http.get<Player[]>(path);
  }

  kick(clanID: number, playerID: number): Observable<any> {
    const path = `${environment.apiUrl}/clans/${clanID}/kick/${playerID}`;
    return this.http.post<any>(path, {});
  }

  leave(): Observable<any> {
    const path = `${environment.apiUrl}/clans/leave`;
    return this.http.post<any>(path, {});
  }
}
