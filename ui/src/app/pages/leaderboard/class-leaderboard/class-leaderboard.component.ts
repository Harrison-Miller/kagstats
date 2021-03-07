import { Component, OnDestroy, OnInit } from '@angular/core';
import { LeaderboardService } from '../../../services/leaderboard.service';
import { Subject } from 'rxjs';
import {BasicStats, PlayerClaims} from '../../../models';
import { takeUntil } from 'rxjs/operators';
import { ActivatedRoute, Router } from '@angular/router';
import {AuthService} from "../../../services/auth.service";

const availableClasses = ['Archer', 'Builder', 'Knight', 'Kills', 'MonthlyArcher', 'MonthlyBuilder', 'MonthlyKnight'];
const boardTitle = {
  'Archer': 'All Time Archer',
  'Builder': 'All Time Builder',
  'Knight': 'All Time Knight',
  'MonthlyArcher': 'Monthly Archer',
  'MonthlyBuilder': 'Monthly Builder',
  'MonthlyKnight': 'Monthly Knight',
  'Kills': 'Kills'
}

const monthNames =["January", "February", "March", "April", "May", "June", "July", "August", "September", "November", "December"];

@Component({
  selector: 'app-class-leaderboard',
  templateUrl: './class-leaderboard.component.html',
  styleUrls: ['./class-leaderboard.component.sass']
})
export class ClassLeaderboardComponent implements OnInit, OnDestroy {
  class: string;
  componentDestroyed$ = new Subject();
  leaderboard$ = new Subject<BasicStats[]>();

  playerClaims: PlayerClaims = null;

  constructor(
    private leaderboardService: LeaderboardService,
    private route: ActivatedRoute,
    private router: Router,
    private authService: AuthService
  ) {}

  ngOnInit() {
    this.authService.playerClaims.pipe(takeUntil(this.componentDestroyed$))
      .subscribe( value => {
        this.playerClaims = value;
    });

    this.route.paramMap
      .pipe(takeUntil(this.componentDestroyed$))
      .subscribe(params => {
        this.class = params.get('board');

        if (!availableClasses.includes(this.class)) {
          this.router.navigateByUrl('/MonthlyArcher');
          console.error('Invalid class');
        }

        this.leaderboardService
          .leaderboard
          .pipe(takeUntil(this.componentDestroyed$))
          .subscribe(result => {
            this.leaderboard$.next(result.leaderboard);
          });

        this.leaderboardService.getLeaderboard(this.class);
      });
  }

  ngOnDestroy() {
    this.componentDestroyed$.next();
  }

  totalKills(leader: BasicStats): number {
    if(this.class == "Kills") {
      return leader.totalKills;
    }
    return leader[`${this.class.toLowerCase().replace("monthly", "")}Kills`]
  }

  boardTitle(): string {
    return boardTitle[this.class];
  }

  boardDate(): string {
    if(this.class.toLocaleLowerCase().includes("monthly")) {
      var date = new Date();
      var year = date.getFullYear()
      var month = monthNames[date.getMonth()];
      return `- ${month} ${year}`
    }
    return "";
  }

  totalDeaths(leader: BasicStats): number {
    if(this.class == "Kills") {
      return leader.totalDeaths;
    }
    return leader[`${this.class.toLowerCase().replace("monthly", "")}Deaths`]
  }

  kd(leader: BasicStats): string {
    const deaths = this.totalDeaths(leader);
    return (this.totalKills(leader) / (deaths === 0 ? 1 : deaths)).toFixed(2);
  }
}
