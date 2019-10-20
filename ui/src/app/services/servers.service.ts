import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Server, APIServer } from '../models';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { environment } from '../../environments/environment';

@Injectable({
  providedIn: 'root'
})
export class ServersService {
  constructor(private http: HttpClient) {}

  getServers(): Observable<Server[]> {
    return this.http.get<Server[]>(`${environment.apiUrl}/servers`);
  }

  getServer(id: number): Observable<Server> {
    return this.http.get<Server>(`${environment.apiUrl}/servers/${id}`);
  }

  getAPIServer(address: string, port: string): Observable<APIServer> {
    let path = `https://api.kag2d.com/server/ip/${address}/port/${port}/status`
    return this.http.get<{serverStatus:APIServer}>(path).pipe(
      map(status => status.serverStatus)
    );
  }
}
