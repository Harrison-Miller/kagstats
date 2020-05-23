import { Component, OnInit } from '@angular/core';
import { MapBasics, GithubTree } from '../../models';
import { MapsService } from '../../services/maps.service';

@Component({
  selector: 'app-maps',
  templateUrl: './maps.component.html',
  styleUrls: ['./maps.component.scss']
})
export class MapsComponent implements OnInit {
  maps: MapBasics[];
  paths: GithubTree

  constructor(private mapsService: MapsService) { }

  ngOnInit() {
    this.mapsService.getMaps().subscribe(
      maps => {
        this.maps = maps;
        this.maps.sort((a,b) => {
          return b.matches - a.matches;
        });
      this.mapsService.getMapPaths().subscribe(
        paths => {
          this.paths = paths;
          for(var map of maps) {
            for(var p of this.paths.tree) {
              if(p.path.endsWith(map.mapName + ".png")) {
                map.image = "https://raw.githubusercontent.com/transhumandesign/kag-base/master/" + p.path;
              }
            }
          }
      });
    });
  }

}
