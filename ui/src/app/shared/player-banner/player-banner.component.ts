import { Component, OnInit, Input } from '@angular/core';
import { Player } from '../../models';
import {Router} from "@angular/router";

export const REGISTERED_TOOLTIP = {
  "r10yBadge": "playing for over 10 years",
  "r9yBadge": "playing for over 9 years",
  "r8yBadge": "playing for over 8 years",
  "r7yBadge": "playing for over 7 years",
  "r6yBadge": "playing for over 6 years",
  "r5yBadge": "playing for over 5 years",
  "r4yBadge": "playing for over 4 years",
  "r3yBadge": "playing for over 3 years",
  "r2yBadge": "playing for over 2 years",
  "r1yBadge": "playing for 1 year",
  "r9mBadge": "playing for 9 months",
  "r6mBadge": "playing for 6 months",
  "r3mBadge": "playing for 3 months",
  "r2mBadge": "playing for 2 months",
  "r1mBadge": "playing for 1 month",
  "r3wBadge": "playing for 3 weeks",
  "r2wBadge": "playing for 2 weeks",
  "r1wBadge": "playing for 1 week",
  "newBadge": "welcome this player to the game!"
};

@Component({
  selector: 'app-player-banner',
  templateUrl: './player-banner.component.html',
  styleUrls: ['./player-banner.component.scss']
})
export class PlayerBannerComponent implements OnInit {

  public registered_tooltip = REGISTERED_TOOLTIP

  @Input() player: Player;
  today: number = Date.now();

  // optional inputs for changing the style of the player banner
  @Input() avatar: boolean = true; // show the avatar
  @Input() status: boolean =  true; // show the status icon
  @Input() accolades: boolean = true; // show accolade list

  similarClanTag: boolean = false;

  constructor(private router: Router) {

  }

  ngOnInit() {
    if (this.player.clanID && this.player.clantag !== '' && this.player.clantag.toLowerCase().includes(this.player.clanInfo.name.toLowerCase())) {
      this.similarClanTag = true;
    }
  }

}
