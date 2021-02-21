import { Component, OnInit } from '@angular/core';
import {ClanInfo} from "../../models";
import {ClansService} from "../../services/clans.service";

@Component({
  selector: 'app-clans',
  templateUrl: './clans.component.html',
  styleUrls: ['./clans.component.sass']
})
export class ClansComponent implements OnInit {

  clans: ClanInfo[];

  constructor(
    private clansService: ClansService
  ) { }

  ngOnInit() {
    this.clansService.getClans().subscribe( resp => {
      this.clans = resp;
    });
  }

}
