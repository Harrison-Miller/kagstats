import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ArcherLeaderboardComponent } from './archer-leaderboard.component';

describe('ArcherLeaderboardComponent', () => {
  let component: ArcherLeaderboardComponent;
  let fixture: ComponentFixture<ArcherLeaderboardComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ArcherLeaderboardComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ArcherLeaderboardComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
