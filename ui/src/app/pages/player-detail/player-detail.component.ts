import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { PlayersService } from '../../services/players.service';
import { Player, BasicStats } from '../../models';

@Component({
  selector: 'app-player-detail',
  templateUrl: './player-detail.component.html',
  styleUrls: ['./player-detail.component.scss']
})
export class PlayerDetailComponent implements OnInit {

  playerId: number;
  player: Player;
  basicStats: BasicStats;

  constructor(
    private route: ActivatedRoute,
    private playersService: PlayersService) { }

  ngOnInit() {
    this.playerId = +this.route.snapshot.paramMap.get('id')

    this.getPlayer();
  }

  getPlayer(): void {
    this.playersService.getPlayer(this.playerId)
      .subscribe( b => {
        this.basicStats = b;
        this.player = this.basicStats.player;
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
