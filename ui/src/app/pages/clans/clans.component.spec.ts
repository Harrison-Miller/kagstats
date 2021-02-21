import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ClansComponent } from './clans.component';

describe('ClansComponent', () => {
  let component: ClansComponent;
  let fixture: ComponentFixture<ClansComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ClansComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ClansComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
