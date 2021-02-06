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
          a.percent = Math.floor(a.wins / a.ballots * 100)
          b.percent = Math.floor(b.wins / b.ballots * 100)
          return b.percent - a.percent;
        });
      this.mapsService.getMapPaths().subscribe(
        paths => {
          this.paths = paths;
          for(var map of maps) {
            for(var p of this.paths.tree) {
              var path = p.path.toLowerCase()
              if(path.endsWith(map.mapName.toLowerCase() + ".png")) {
                map.gamemode = "other";
                if (path.includes("smallctf")) {
                  map.gamemode = "Small CTF";
                }
                else if (path.includes("ctf")) {
                  map.gamemode = "CTF";
                } else if (path.includes("tth") || path.includes("war")) {
                  map.gamemode = "TTH";
                }
                console.log(path);
                map.image = "https://raw.githubusercontent.com/transhumandesign/kag-base/master/" + p.path;
              }
            }
          }
      });
    });
  }

}
