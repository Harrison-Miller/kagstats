import { Component, OnInit } from '@angular/core';
import {PollService} from "../../services/poll.service";
import {Poll} from "../../models";

@Component({
  selector: 'app-survey-management',
  templateUrl: './survey-management.component.html',
  styleUrls: ['./survey-management.component.sass']
})
export class SurveyManagementComponent implements OnInit {
  polls: Poll[];

  constructor(private pollService: PollService) { }

  ngOnInit() {
    this.pollService.listPolls().subscribe(polls => {
      this.polls = polls;
    });
  }

  downloadResponses(pollID: number): void {
    console.log('download responses for ' + pollID);
    this.pollService.downloadPoll(pollID).subscribe(resp => {
      let fileName = resp.headers.get('Content-Disposition');
      fileName = fileName.split('=')[1];
      console.log('fileName: ', fileName);

      // It is necessary to create a new blob object with mime-type explicitly set
      // otherwise only Chrome works like it should
      const newBlob = new Blob([resp.body], { type: 'text/csv' });

      // IE doesn't allow using a blob object directly as link href
      // instead it is necessary to use msSaveOrOpenBlob
      if (window.navigator && window.navigator.msSaveOrOpenBlob) {
        window.navigator.msSaveOrOpenBlob(newBlob, fileName);
        return;
      }

      // For other browsers:
      // Create a link pointing to the ObjectURL containing the blob.
      const data = window.URL.createObjectURL(newBlob);

      const link = document.createElement('a');
      link.href = data;
      link.download = fileName;
      // this is necessary as link.click() does not work on the latest firefox
      link.dispatchEvent(new MouseEvent('click', { bubbles: true, cancelable: true, view: window }));

      setTimeout(function () {
        // For Firefox it is necessary to delay revoking the ObjectURL
        window.URL.revokeObjectURL(data);
        link.remove();
      }, 100);
    }, error => {
      console.log(error);
    });
  }
}
