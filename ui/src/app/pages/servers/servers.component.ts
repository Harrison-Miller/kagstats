import { Component, OnInit } from '@angular/core';
import { Server } from '../../models';
import { ServersService } from '../../services/servers.service';

@Component({
  selector: 'app-servers',
  templateUrl: './servers.component.html',
  styleUrls: ['./servers.component.sass']
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
      .subscribe(servers => (this.servers = servers));
  }
}
