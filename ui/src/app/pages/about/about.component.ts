import { Component, OnInit } from '@angular/core';


export const FAQ: {question: string, answer: string}[] = [
  { 
    question: 'How can I report a bug or request a feature?',
    answer: 'Go to the <a href="https://github.com/Harrison-Miller/kagstats">github</a> repository and submit an issue or feature request. You can also reach out and discuss KAG stats on the community <a href="https://discordapp.com/invite/kag">discord</a>.'
  },
  {
    question: 'How can I host my own KAG stats website?',
    answer: ' If you want to host your own KAG stats website great! Visit the documentation in the github repo (docs yet to be written, link TBD) for general hosting instructions. TL;DR you will need host that supports docker, docker-compose.'
  },
  {
    question: 'How often are the stats updated?',
    answer: "It'll take about 2 minutes before you see the most recent information in the system. The collector (connected to the game server) adds entries to the database about every minute."
  },
  {
    question: "My K/D is high, why can't I see my self on the leaderboard?",
    answer: 'There are measures in place to prevent unfair placements on the leaderboards. You must have enough entries in the database and time played before you can potentially appear on the leaderboard. These measures may be adjusted or added to as the system improves.',
  }
];


@Component({
  selector: 'app-about',
  templateUrl: './about.component.html',
  styleUrls: ['./about.component.scss']
})
export class AboutComponent implements OnInit {

  public faq = FAQ;

  constructor() { }

  ngOnInit() {
  }

}
