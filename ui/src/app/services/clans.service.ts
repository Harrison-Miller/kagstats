import { Injectable } from '@angular/core';
import {Observable} from "rxjs";
import {HttpClient} from "@angular/common/http";
import { environment } from '../../environments/environment';
import {ClanInfo} from "../models";

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
}
