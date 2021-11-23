import { Component, OnInit } from '@angular/core';
import {PollService} from "../../services/poll.service";
import {Poll} from "../../models";

@Component({
  selector: 'app-poll',
  templateUrl: './poll.component.html',
  styleUrls: ['./poll.component.sass']
})
export class PollComponent implements OnInit {
  poll: Poll;

  constructor(private pollService: PollService) { }

  ngOnInit() {
    this.pollService.getCurrentPoll().subscribe(
      poll => {
        this.poll = poll;
        for(var q of this.poll.questions) {
          q.options_split = q.options.split(';');
        }
      }
    );
  }

}
