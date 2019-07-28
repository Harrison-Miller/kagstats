import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { ReactiveFormsModule } from '@angular/forms';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { PlayersService } from './services/players.service';
import { HttpClientModule } from '@angular/common/http';
import { KagTableComponent } from './shared/kag-table/kag-table.component';
import { HomeComponent } from './pages/home/home.component';
import { LeaderboardComponent } from './pages/leaderboard/leaderboard.component';
import { BaseLeaderboardComponent } from './pages/leaderboard/base-leaderboard/base-leaderboard.component';
import { ClassLeaderboardComponent } from './pages/leaderboard/class-leaderboard/class-leaderboard.component';
import { LeaderboardService } from './services/leaderboard.service';
import { ServersComponent } from './pages/servers/servers.component';
import { ServersService } from './services/servers.service';
import { PlayersComponent } from './pages/players/players.component';
import { UnkownAvatarPipe } from './pipes/unkown-avatar.pipe';
import { PlayerBannerComponent } from './shared/player-banner/player-banner.component';
import { StatusClassPipe } from './pipes/status-class.pipe';
import { KagRoleColorPipe } from './pipes/kag-role-color.pipe';
import { DefaultImageDirective } from './default-image.directive';
import { PlayerDetailComponent } from './pages/player-detail/player-detail.component';

@NgModule({
  declarations: [
    AppComponent,
    KagTableComponent,
    HomeComponent,
    LeaderboardComponent,
    BaseLeaderboardComponent,
    ClassLeaderboardComponent,
    ServersComponent,
    PlayersComponent,
    UnkownAvatarPipe,
    PlayerBannerComponent,
    StatusClassPipe,
    KagRoleColorPipe,
    DefaultImageDirective,
    PlayerDetailComponent,
  ],
  imports: [BrowserModule, AppRoutingModule, HttpClientModule, ReactiveFormsModule],
  providers: [LeaderboardService, PlayersService, ServersService],
  bootstrap: [AppComponent]
})
export class AppModule {}
