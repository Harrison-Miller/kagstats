import { Component, OnInit, ViewChild, OnDestroy } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { PlayersService } from '../../services/players.service';
import { Player, BasicStats, Nemesis, Hitter, Server, APIPlayerStatus } from '../../models';
import { NemesisService } from '../../services/nemesis.service';
import { takeUntil } from 'rxjs/operators';
import { HITTER_DESCRIPTION } from '../../hitters';
import { REGISTERED_TOOLTIP } from '../../shared/player-banner/player-banner.component';
import { HittersService } from '../../services/hitters.service';
import { ServersService } from '../../services/servers.service';
import { Title } from '@angular/platform-browser';
import { CapturesService } from '../../services/captures.service';

@Component({
  selector: 'app-player-detail',
  templateUrl: './player-detail.component.html',
  styleUrls: ['./player-detail.component.scss', '../../shared/killfeed/killfeed.component.scss', '../../shared/player-banner/player-banner.component.scss']
})
export class PlayerDetailComponent implements OnInit, OnDestroy {

  public registered_tooltip = REGISTERED_TOOLTIP;

  playerId: number;
  player: Player;
  basicStats: BasicStats = {
    player: null,
    suicides: 0,
    teamKills: 0,
    archerKills: 0,
    archerDeaths: 0,
    builderKills: 0,
    builderDeaths: 0,
    knightKills: 0,
    knightDeaths: 0,
    otherKills: 0,
    otherDeaths: 0,
    totalKills: 0,
    totalDeaths: 0
  };

  status: APIPlayerStatus;

  captures: number;
  nemesis: Nemesis;
  bullied: Nemesis[];
  hitters: Hitter[];
  descriptions: string[] = HITTER_DESCRIPTION;

  today: number = Date.now();

  @ViewChild('t') t;

  constructor(
    private route: ActivatedRoute,
    private playersService: PlayersService,
    private nemesisService: NemesisService,
    private hittersService: HittersService,
    private serversService: ServersService,
    private titleService: Title,
    private capturesService: CapturesService) { }

  ngOnInit() {
    this.route.paramMap.subscribe(params => {
      this.playerId = +params.get('id');
      window.scrollTo(0, 0);

      this.nemesis = null;
      this.bullied = null;
      this.hitters = null;
      this.captures = 0;


      if(this.t) {
        this.t.select('overview');
      }
      
      this.getPlayer();
      this.nemesisService.getNemesis(this.playerId)
        .subscribe( nemesis => {
          if(nemesis) {
            this.nemesis = nemesis;
          }
        });
      this.nemesisService.getBullied(this.playerId)
        .subscribe( bullied => {
          if(bullied) {
            
            this.bullied = bullied.sort((a,b) => b.deaths - a.deaths);
          }
        });
      this.hittersService.getHitters(this.playerId)
        .subscribe( hitters => {
          if(hitters) {
            this.hitters = hitters.slice(0, 3);
          }
        });

      this.capturesService.getCaptures(this.playerId)
        .subscribe( c => {
          if(c) {
            this.captures = c.captures;
          }
        })
    });
  }

  ngOnDestroy() {
    this.titleService.setTitle("KAG Stats");
  }

  getPlayer(): void {
    this.playersService.getPlayer(this.playerId)
      .subscribe( b => {
        this.basicStats = b;
        this.player = this.basicStats.player;
        this.getStatus(this.player.username);
        this.titleService.setTitle("KAG Stats - " + this.player.characterName);
      },
      error => {
        this.playersService.getPlayerName(this.playerId)
          .subscribe( a => {
            this.player = a;
            this.getStatus(this.player.username);
            this.titleService.setTitle("KAG Stats - " + this.player.characterName);
          })
      });
  }

  getStatus(username: string): void {
    this.playersService.getStatus(username)
      .subscribe( status => {
        this.status = status;
        this.getServer(status.server.serverIPv4Address, status.server.serverPort);
      });
  }

  getServer(ip: string, port: string): void {
    this.serversService.getAPIServer(ip,  port)
      .subscribe( server => {
        this.status.apiServer = server;
      });
  }

  totalKD(): string {
    let kills = this.basicStats.totalKills;
    let deaths = this.basicStats.totalDeaths;
    return (kills / (deaths === 0 ? 1 : deaths)).toFixed(2);
  }

  archerKD(): string {
    let kills = this.basicStats.archerKills;
    let deaths = this.basicStats.archerDeaths;
    return (kills / (deaths === 0 ? 1 : deaths)).toFixed(2);
  }

  builderKD(): string {
    let kills = this.basicStats.builderKills;
    let deaths = this.basicStats.builderDeaths;
    return (kills / (deaths === 0 ? 1 : deaths)).toFixed(2);
  }

  knightKD(): string {
    let kills = this.basicStats.knightKills;
    let deaths = this.basicStats.knightDeaths;
    return (kills / (deaths === 0 ? 1 : deaths)).toFixed(2);
  }

}
