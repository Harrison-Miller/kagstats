import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { HomeComponent } from './pages/home/home.component';
import { ServersComponent } from './pages/servers/servers.component';
import { LeaderboardComponent } from './pages/leaderboard/leaderboard.component';
import { BaseLeaderboardComponent } from './pages/leaderboard/base-leaderboard/base-leaderboard.component';
import { ClassLeaderboardComponent } from './pages/leaderboard/class-leaderboard/class-leaderboard.component';
import { PlayersComponent } from './pages/players/players.component';
import { PlayerDetailComponent } from './pages/player-detail/player-detail.component';
import { KillsComponent } from './pages/kills/kills.component';
import { AboutComponent } from './pages/about/about.component';

const routes: Routes = [{
  path: '',
  redirectTo: 'leaderboards',
  pathMatch: 'full'
}, {
  path: 'leaderboards',
  component: LeaderboardComponent,
  children: [{
    path: '',
    component: BaseLeaderboardComponent
  }, {
    path: ':board',
    component: ClassLeaderboardComponent
  }]
}, {
  path: 'servers',
  component: ServersComponent
},{
  path: 'players',
  component: PlayersComponent
},{
  path: 'players/:id',
  component: PlayerDetailComponent
},{
  path: 'kills',
  component: KillsComponent
},{
  path: 'about',
  component: AboutComponent
}];

@NgModule({
  imports: [RouterModule.forRoot(routes, { useHash: true })],
  exports: [RouterModule]
})
export class AppRoutingModule { }
