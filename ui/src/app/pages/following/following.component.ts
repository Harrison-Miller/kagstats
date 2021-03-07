import {Component, OnDestroy, OnInit} from '@angular/core';
import {AuthService} from '../../services/auth.service';
import {PlayersService} from '../../services/players.service';
import {FollowersService} from '../../services/followers.service';
import {takeUntil} from 'rxjs/operators';
import {BasicStats, PlayerClaims} from '../../models';
import {Subject} from 'rxjs';
import {ServersService} from '../../services/servers.service';

@Component({
  selector: 'app-following',
  templateUrl: './following.component.html',
  styleUrls: ['./following.component.sass']
})
export class FollowingComponent implements OnInit, OnDestroy {

  playerClaims: PlayerClaims = null;
  componentDestroyed$ = new Subject();

  following: BasicStats[];
  leaderboardClass = 'Archer';
  leaderboard = [];
  archer = [];
  builder = [];
  knight = [];


  constructor(private authService: AuthService,
              private playersService: PlayersService,
              private followersService: FollowersService,
              private serversService: ServersService) { }

  ngOnInit() {
    this.authService.playerClaims.pipe(takeUntil(this.componentDestroyed$))
      .subscribe( value => {
        this.playerClaims = value;
        if (this.playerClaims) {
          this.followersService.getFollowedStats().subscribe( res => {
            this.following = res;
            this.playersService.getPlayer(this.playerClaims.playerID).subscribe(res => {
              this.following.push(res);
              this.setupLeaderboards();

              // query status of everyone
              for (const s of this.following) {
                this.playersService.getStatus(s.player.username).subscribe( status => {
                  this.serversService.getAPIServer(status.server.serverIPv4Address,  status.server.serverPort)
                    .subscribe( server => {
                      status.apiServer = server;
                      this.updateStatus(s.player.username, status);
                    });
                });
              }
            });
          });
        }
      });
  }

  ngOnDestroy(): void {
    this.componentDestroyed$.next();
  }

  getKills(stats): number {
    if (this.leaderboardClass === 'Builder') {
      return stats.builderKills;
    } else if (this.leaderboardClass === 'Knight') {
      return stats.knightKills;
    }
    return stats.archerKills;
  }

  getDeaths(stats): number {
    if (this.leaderboardClass === 'Builder') {
      return stats.builderDeaths;
    } else if (this.leaderboardClass === 'Knight') {
      return stats.knightDeaths;
    }
    return stats.archerDeaths;
  }

  getLeaderboardKDR(stats): string {
    return this.getKDR(this.getKills(stats), this.getDeaths(stats));
  }

  getKDR(kills: number, deaths: number): string {
    return (kills / (deaths === 0 ? 1 : deaths)).toFixed(2);
  }

  setupLeaderboards() {
    this.archer = [...this.following];
    this.archer = this.archer.sort((a, b) => {
      const aKDR = parseFloat(this.getKDR(a.archerKills, a.archerDeaths));
      const bKDR = parseFloat(this.getKDR(b.archerKills, b.archerDeaths));
      return bKDR - aKDR;
    });

    this.builder = [...this.following];
    this.builder = this.builder.sort((a, b) => {
      const aKDR = parseFloat(this.getKDR(a.builderKills, a.builderDeaths));
      const bKDR = parseFloat(this.getKDR(b.builderKills, b.builderDeaths));
      return bKDR - aKDR;
    });

    this.knight = [...this.following];
    this.knight = this.knight.sort((a, b) => {
      const aKDR = parseFloat(this.getKDR(a.knightKills, a.knightDeaths));
      const bKDR = parseFloat(this.getKDR(b.knightKills, b.knightDeaths));
      return bKDR - aKDR;
    });

    this.archerLeaderboard();
  }

  updateStatus(username, status) {
    for (const i in this.archer) {
      if (this.archer[i].player.username === username) {
        this.archer[i].status = status;
      }
    }

    for (const i in this.builder) {
      if (this.builder[i].player.username === username) {
        this.builder[i].status = status;
      }
    }

    for (const i in this.knight) {
      if (this.knight[i].player.username === username) {
        this.knight[i].status = status;
      }
    }
  }

  archerLeaderboard() {
    this.leaderboardClass = 'Archer';
    this.leaderboard = this.archer;
  }

  builderLeaderboard() {
    this.leaderboardClass = 'Builder';
    this.leaderboard = this.builder;
  }

  knightLeaderboard() {
    this.leaderboardClass = 'Knight';
    this.leaderboard = this.knight;
  }

}
