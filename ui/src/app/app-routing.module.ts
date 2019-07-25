import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { HomeComponent } from './pages/home/home.component';
import { LeaderboardComponent } from './pages/leaderboard/leaderboard.component';
import { BaseLeaderboardComponent } from './pages/leaderboard/base-leaderboard/base-leaderboard.component';
import { ClassLeaderboardComponent } from './pages/leaderboard/class-leaderboard/class-leaderboard.component';

const routes: Routes = [{
  path: '',
  pathMatch: 'full',
  component: HomeComponent
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
}];

@NgModule({
  imports: [RouterModule.forRoot(routes, { useHash: true })],
  exports: [RouterModule]
})
export class AppRoutingModule { }
