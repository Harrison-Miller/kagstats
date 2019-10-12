import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { ReactiveFormsModule } from '@angular/forms';
import { NgbModule } from '@ng-bootstrap/ng-bootstrap';

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
import { KillsService } from './services/kills.service';
import { KagClassImgPipe } from './pipes/kag-class-img.pipe';
import { TimeAgoPipe } from 'time-ago-pipe';
import { HitterIconPipe } from './pipes/hitter-icon.pipe';
import { KillsComponent } from './pages/kills/kills.component';
import { KillfeedComponent } from './shared/killfeed/killfeed.component';
import { AboutComponent } from './pages/about/about.component';
import { StatusService } from './services/status.service';
import { RegisteredAccoladePipe } from './pipes/registered-accolade.pipe';

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
    KillfeedComponent,
    KagClassImgPipe,
    TimeAgoPipe,
    HitterIconPipe,
    KillsComponent,
    AboutComponent,
    RegisteredAccoladePipe
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    HttpClientModule,
    ReactiveFormsModule,
    NgbModule
  ],
  providers: [LeaderboardService, PlayersService, ServersService, KillsService, StatusService],
  bootstrap: [AppComponent]
})
export class AppModule {}
