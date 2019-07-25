import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ClassLeaderboardComponent } from './class-leaderboard.component';

describe('ClassLeaderboardComponent', () => {
  let component: ClassLeaderboardComponent;
  let fixture: ComponentFixture<ClassLeaderboardComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ClassLeaderboardComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ClassLeaderboardComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
