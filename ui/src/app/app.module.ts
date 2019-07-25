import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { PlayersService } from './services/players.service';
import { HttpClientModule } from '@angular/common/http';
import { KagTableComponent } from './shared/kag-table/kag-table.component';
import { HomeComponent } from './pages/home/home.component';
import { LeaderboardComponent } from './pages/leaderboard/leaderboard.component';
import { BaseLeaderboardComponent } from './pages/leaderboard/base-leaderboard/base-leaderboard.component';
import { ArcherLeaderboardComponent } from './pages/leaderboard/archer-leaderboard/archer-leaderboard.component';
import { LeaderboardService } from './services/leaderboard.service';

@NgModule({
  declarations: [
    AppComponent,
    KagTableComponent,
    HomeComponent,
    LeaderboardComponent,
    BaseLeaderboardComponent,
    ArcherLeaderboardComponent
  ],
  imports: [BrowserModule, AppRoutingModule, HttpClientModule],
  providers: [LeaderboardService, PlayersService],
  bootstrap: [AppComponent]
})
export class AppModule {}