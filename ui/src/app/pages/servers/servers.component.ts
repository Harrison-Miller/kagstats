import { Component, OnInit } from '@angular/core';
import { Server } from '../../models';
import { ServersService } from '../../services/servers.service';

@Component({
  selector: 'app-servers',
  templateUrl: './servers.component.html',
  styleUrls: ['./servers.component.scss']
})
export class ServersComponent implements OnInit {
  servers: Server[];

  constructor(private serversService: ServersService) {
    this.getServers();
  }

  ngOnInit() {}

  getServers(): void {
    this.serversService
      .getServers()
      .subscribe(servers => {
        this.servers = servers;
        this.getAPIServers();
      });
  }

  getAPIServers(): void {
    this.servers.forEach(server => {
      this.serversService.getAPIServer(server.address, server.port)
        .subscribe(status => {
          server.APIStatus = status;
          this.servers.sort((a,b) => {
            if (!a.status) {
              return 1;
            }else if(!b.status) {
              return -1;
            }
            return b.APIStatus.currentPlayers - a.APIStatus.currentPlayers;
          });
        })
    });
  }
}
