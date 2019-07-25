import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { PlayersService } from './services/players.service';
import { HttpClientModule } from '@angular/common/http';
import { KagTableComponent } from './shared/kag-table/kag-table.component';
import { HomeComponent } from './pages/home/home.component';
import { LeaderboardComponent } from './pages/leaderboard/leaderboard.component';
import { ServersComponent } from './pages/servers/servers.component';

@NgModule({
  declarations: [AppComponent, KagTableComponent, HomeComponent, LeaderboardComponent, ServersComponent],
  imports: [BrowserModule, AppRoutingModule, HttpClientModule],
  providers: [PlayersService],
  bootstrap: [AppComponent]
})
export class AppModule {}
