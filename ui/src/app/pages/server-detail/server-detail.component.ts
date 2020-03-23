import { Component, OnInit } from '@angular/core';
import { Server, Player } from '../../models';
import { ActivatedRoute } from '@angular/router';
import { ServersService } from 'src/app/services/servers.service';
import { PlayersService } from 'src/app/services/players.service';

@Component({
  selector: 'app-server-detail',
  templateUrl: './server-detail.component.html',
  styleUrls: ['./server-detail.component.sass']
})
export class ServerDetailComponent implements OnInit {

  constructor(
    private route: ActivatedRoute,
    private serversService: ServersService,
    private playersService: PlayersService
  ) { }

  serverId: number;
  server: Server;
  currentPlayers: Player[];

  ngOnInit() {
    this.route.paramMap.subscribe(params => {
      this.serverId = +params.get('id');
    });

    // Retrieve component info.
    this.getServer(this.serverId);

    // Init to empty array, in order to push to later.
    this.currentPlayers = [];
  }

  /**
   * Given the server id passed to the component,
   * retrieve related API info about the server from both
   * the stats API and the kag2d API.
   * 
   * @param serverId The stat's API's own id for the given server.
   */
  getServer(serverId: number): void{
    this.serversService.getServer(this.serverId)
      .subscribe( s => {
        this.server = s;
        // Fill our model's APIStatus object.
        this.getApiServer();
      },
      error => {})
  }

  /**
   * Once we've retrieved the kagstats api model, perform an api call
   * on the kag2d api, and update our current server model instance.
   */
  getApiServer(): void {
    this.serversService.getAPIServer(this.server.address, this.server.port)
      .subscribe(status => {
        
        // Update model's field.
        this.server.APIStatus = status;

        // From the simple playerList given by kag2d api, retrieve player model data from kagstats api.
        // TODO: This has gotta be costly, can we perform a bulk request/lookup?
        status.playerList.forEach(playerJson => {
          if(playerJson)
            this.playersService.searchPlayers(playerJson["username"])
              .subscribe(players => this.currentPlayers.push(players[0]));
        });
      })
  }
}
