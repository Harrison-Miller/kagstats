import { Component, OnInit, ViewChild, OnDestroy } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { PlayersService } from '../../services/players.service';
import { Player, BasicStats, Nemesis, Hitter } from '../../models';
import { NemesisService } from '../../services/nemesis.service';
import { takeUntil } from 'rxjs/operators';
import { HITTER_DESCRIPTION } from '../../hitters';
import { HittersService } from '../../services/hitters.service';
import { Title } from '@angular/platform-browser';

@Component({
  selector: 'app-player-detail',
  templateUrl: './player-detail.component.html',
  styleUrls: ['./player-detail.component.scss', '../../shared/killfeed/killfeed.component.scss']
})
export class PlayerDetailComponent implements OnInit, OnDestroy {

  playerId: number;
  player: Player;
  basicStats: BasicStats;
  nemesis: Nemesis;
  hitters: Hitter[];
  descriptions: string[] = HITTER_DESCRIPTION;

  @ViewChild('t') t;

  constructor(
    private route: ActivatedRoute,
    private playersService: PlayersService,
    private nemesisService: NemesisService,
    private hittersService: HittersService,
    private titleService: Title) { }

  ngOnInit() {
    this.route.paramMap.subscribe(params => {
      this.playerId = +params.get('id');
      window.scrollTo(0, 0);

      if(this.t) {
        this.t.select('overview');
      }
      
      this.getPlayer();
      this.nemesisService.getNemesis(this.playerId)
        .subscribe( nemeses => {
          this.nemesis = nemeses[0];
        });
      this.hittersService.getHitters(this.playerId)
        .subscribe( hitters => {
          this.hitters = hitters.slice(0, 3);
        });
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
        this.titleService.setTitle("KAG Stats - " + this.player.characterName);
      });
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
