import { Component, OnInit } from '@angular/core';
import {ActivatedRoute} from "@angular/router";
import {ClansService} from "../../services/clans.service";
import {BasicStats, ClanInfo, Player} from "../../models";

@Component({
  selector: 'app-clan-detail',
  templateUrl: './clan-detail.component.html',
  styleUrls: ['./clan-detail.component.scss']
})
export class ClanDetailComponent implements OnInit {

  clanID: number;
  clan: ClanInfo;
  members: BasicStats[];

  constructor(
    private route: ActivatedRoute,
    private clanService: ClansService,
  ) { }

  ngOnInit() {
    this.route.paramMap.subscribe( params => {
      this.clanID = +params.get('id');
      window.scrollTo(0, 0);
      this.updateClan();

    });
  }

  updateClan(): void {
    this.clanService.getClan(this.clanID).subscribe( resp => {
      this.clan = resp;
    });

    this.clanService.getMembers(this.clanID).subscribe( resp => {
      this.members = resp;
    });
  }

  getKD(kills: number, deaths: number): string {
    return (kills / (deaths === 0 ? 1 : deaths)).toFixed(2);
  }

}
