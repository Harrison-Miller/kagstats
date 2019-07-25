import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Server } from '../models';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})
export class ServersService {
  constructor(private http: HttpClient) {}

  getServers(): Observable<Server[]> {
    return this.http.get<{ servers: Server[] }>('/api/servers').pipe(
      map(response => {
        return response.servers;
      })
    );
  }
}
