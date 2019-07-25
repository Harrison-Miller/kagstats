import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { BaseLeaderboardComponent } from './base-leaderboard.component';

describe('BaseLeaderboardComponent', () => {
  let component: BaseLeaderboardComponent;
  let fixture: ComponentFixture<BaseLeaderboardComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ BaseLeaderboardComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(BaseLeaderboardComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
