import {Component, OnDestroy, OnInit} from '@angular/core';
import {PollService} from "../../services/poll.service";
import {PlayerClaims, Poll, PollAnswer} from "../../models";
import {AuthService} from "../../services/auth.service";
import {takeUntil} from "rxjs/operators";
import {Subject} from "rxjs";
import {FormBuilder, FormControl, Validators} from "@angular/forms";
import {Md5} from "ts-md5";

@Component({
  selector: 'app-poll',
  templateUrl: './poll.component.html',
  styleUrls: ['./poll.component.sass']
})
export class PollComponent implements OnInit, OnDestroy {
  poll: Poll;
  pollCompleted: boolean;
  playerClaims: PlayerClaims;
  componentDestroyed$ = new Subject();

  pollForm = this.formBuilder.group({});
  formErrors = {};

  constructor(private pollService: PollService,
              private authService: AuthService,
              private formBuilder: FormBuilder) { }

  getPoll() {
    this.pollService.pollCompleted().subscribe( resp => { this.pollCompleted = resp.completed; });
    this.pollService.getCurrentPoll().subscribe(
      poll => {
        this.poll = poll;
        const controls = {};
        for(var q of this.poll.questions) {
          q.options_split = q.options.split(';');
          q.questionHash = Md5.hashStr(q.question);
          controls[q.questionHash] = '';
          if (q.options_split.includes('Other')) {
            controls[q.questionHash + 'other'] = '';
          }
        }
        this.pollForm = this.formBuilder.group(controls);
      }
    );
  }

  ngOnInit() {
    this.getPoll();

    this.authService.playerClaims.pipe(takeUntil(this.componentDestroyed$))
      .subscribe( value => {
        this.playerClaims = value;
      });
  }

  ngOnDestroy() {
    this.componentDestroyed$.next();

  }

  validate() {
    this.formErrors = {};
    for(var q of this.poll.questions) {
      let value = this.pollForm.value[q.questionHash].trim();
      let isOther = false;
      if (value.includes('Other')) {
        value = this.pollForm.value[q.questionHash + 'other'].trim();
        isOther = true;
      }

      if (q.required && value === '') {
        this.formErrors[q.questionHash] = 'required question';
        continue;
      }

      if (!isOther && q.options_split.length > 1 && !q.options_split.includes(value)) {
        this.formErrors[q.questionHash] = 'choose a valid option';
        continue;
      }
    }

    return Object.keys(this.formErrors).length
  }

  onPollSubmit() {
    const errors = this.validate();
    if (errors > 0) {
      return;
    }

    const answers: PollAnswer[] = [];
    for(var q of this.poll.questions) {
      let value = this.pollForm.value[q.questionHash].trim();
      if (value.includes('Other')) {
        value = this.pollForm.value[q.questionHash + 'other'].trim();
      }

      const a = {} as PollAnswer;
      a.answer = value;
      a.questionID = q.questionID;
      answers.push(a);
    }
    console.log(answers);
    this.pollService.answerCurrentPoll(answers).subscribe( resp => {
      this.getPoll();
    }, err => {
      console.log(err);
    });
  }

}
